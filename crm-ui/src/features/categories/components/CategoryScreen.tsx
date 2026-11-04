"use client";

import { useMemo, useState } from "react";
import type { Key } from "react";
import { Tree, Typography } from "antd";
import type { TreeDataNode } from "antd";
import { AppPageHeader } from "@/components/ui/AppPageHeader";
import { AppError } from "@/components/ui/AppError";
import { AppEmpty } from "@/components/ui/AppEmpty";
import { AppLoading } from "@/components/ui/AppLoading";
import { SearchInput } from "@/components/shared/SearchInput";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { useCategories } from "../hooks/useCategories";
import type { Category } from "../types/category.types";

interface CategoryTree {
  nodes: TreeDataNode[];
  // Keys to expand: every branch when browsing, only the paths to matches when searching.
  expandedKeys: number[];
}

function buildTree(categories: Category[], search: string): CategoryTree {
  const children = new Map<number | null, Category[]>();
  for (const category of categories) {
    const siblings = children.get(category.parent) ?? [];
    siblings.push(category);
    children.set(category.parent, siblings);
  }

  const term = search.trim().toLowerCase();
  const expandedKeys: number[] = [];

  // Returns the visible nodes under `parent` and whether any of them matched the search.
  // Children of a match stay visible so a matched branch can still be browsed.
  const build = (
    parent: number | null,
    ancestorMatched: boolean,
  ): { nodes: TreeDataNode[]; hasMatch: boolean } => {
    let hasMatch = false;
    const nodes = (children.get(parent) ?? [])
      .sort((a, b) => a.order - b.order)
      .flatMap((category) => {
        const matched =
          term !== "" &&
          (category.title.toLowerCase().includes(term) || category.slug.includes(term));
        const subtree = build(category.id, ancestorMatched || matched);
        if (term !== "" && !ancestorMatched && !matched && !subtree.hasMatch) {
          return [];
        }
        hasMatch ||= matched || subtree.hasMatch;
        if (subtree.nodes.length > 0 && (term === "" || subtree.hasMatch)) {
          expandedKeys.push(category.id);
        }
        return [
          {
            key: category.id,
            title: (
              <span>
                {category.title}{" "}
                <Typography.Text type="secondary">({category.slug})</Typography.Text>
              </span>
            ),
            isLeaf: category.isLeaf,
            children: subtree.nodes,
          },
        ];
      });
    return { nodes, hasMatch };
  };

  return { nodes: build(null, false).nodes, expandedKeys };
}

export function CategoryScreen() {
  const { t } = useI18n();
  const query = useCategories();
  const [search, setSearch] = useState("");
  const [expanded, setExpanded] = useState<Key[] | null>(null);

  const tree = useMemo(() => buildTree(query.data ?? [], search), [query.data, search]);

  const onSearch = (value: string) => {
    setSearch(value);
    setExpanded(null);
  };

  return (
    <div>
      <AppPageHeader title={t("categories")} />
      <div className="app-filters">
        <SearchInput value={search} onChange={onSearch} />
      </div>
      {query.error ? <AppError error={query.error} onRetry={() => void query.refetch()} /> : null}
      {query.isLoading ? (
        <AppLoading />
      ) : tree.nodes.length === 0 ? (
        <AppEmpty />
      ) : (
        <Tree
          showLine
          blockNode
          selectable={false}
          treeData={tree.nodes}
          expandedKeys={expanded ?? tree.expandedKeys}
          onExpand={setExpanded}
        />
      )}
    </div>
  );
}
