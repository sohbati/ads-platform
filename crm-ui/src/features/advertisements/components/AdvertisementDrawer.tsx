"use client";

import { useState } from "react";
import type { DescriptionsProps } from "antd";
import {
  Alert,
  Button,
  Descriptions,
  Drawer,
  Empty,
  Image,
  Space,
  Typography,
} from "antd";
import { AppError } from "@/components/ui/AppError";
import { AppLoading } from "@/components/ui/AppLoading";
import { ConfirmDialog } from "@/components/shared/ConfirmDialog";
import { mediaUrl } from "@/config/publicEnv";
import { can } from "@/lib/auth/permissions";
import { useSession } from "@/lib/auth/SessionProvider";
import { formatDateTime } from "@/lib/formatters/date";
import { formatDigits } from "@/lib/formatters/digits";
import { formatMobile } from "@/lib/formatters/phone";
import { useI18n } from "@/lib/i18n/I18nProvider";
import { Permission } from "@/lib/permissions/codes";
import { AdvertisementStatus } from "../constants/advertisement.constants";
import { useAdvertisement, useReviewMutations } from "../hooks/useAdvertisements";
import type { AdvertisementDetail } from "../types/advertisement.types";
import { AdvertisementStatusTag } from "./AdvertisementStatusTag";
import { formatPrice } from "./formatPrice";
import { RejectDialog } from "./RejectDialog";

interface AdvertisementDrawerProps {
  adId: number | null;
  onClose: () => void;
}

function displayValue(value: unknown): string {
  if (value === null || value === undefined || value === "") {
    return "—";
  }
  if (typeof value === "boolean") {
    return value ? "✓" : "✗";
  }
  if (typeof value === "object") {
    return JSON.stringify(value);
  }
  return String(value);
}

export function AdvertisementDrawer({ adId, onClose }: AdvertisementDrawerProps) {
  const { t, locale } = useI18n();
  const session = useSession();
  const query = useAdvertisement(adId);
  const { approve, reject } = useReviewMutations();
  const [approving, setApproving] = useState(false);
  const [rejecting, setRejecting] = useState(false);
  const ad = query.data;

  const pending = ad?.status === AdvertisementStatus.pending;
  const canApprove = pending && can(session.permissions, Permission.advertisementApprove);
  const canReject = pending && can(session.permissions, Permission.advertisementReject);

  return (
    <Drawer
      open={adId !== null}
      title={ad ? ad.title : t("adDetails")}
      size="large"
      onClose={onClose}
      extra={
        canApprove || canReject ? (
          <Space>
            {canReject ? (
              <Button danger onClick={() => setRejecting(true)}>
                {t("reject")}
              </Button>
            ) : null}
            {canApprove ? (
              <Button type="primary" onClick={() => setApproving(true)}>
                {t("approve")}
              </Button>
            ) : null}
          </Space>
        ) : null
      }
    >
      {query.isLoading ? <AppLoading /> : null}
      {query.error ? (
        <AppError error={query.error} onRetry={() => void query.refetch()} />
      ) : null}
      {ad ? <Details ad={ad} locale={locale} t={t} /> : null}

      <ConfirmDialog
        open={approving}
        title={t("confirmApprove")}
        loading={approve.isPending}
        onCancel={() => setApproving(false)}
        onConfirm={() => {
          if (ad) {
            approve.mutate(ad.id, { onSettled: () => setApproving(false) });
          }
        }}
      />
      <RejectDialog
        open={rejecting}
        loading={reject.isPending}
        onCancel={() => setRejecting(false)}
        onSubmit={(input) => {
          if (ad) {
            reject.mutate({ id: ad.id, input }, { onSettled: () => setRejecting(false) });
          }
        }}
      />
    </Drawer>
  );
}

function Details({
  ad,
  locale,
  t,
}: {
  ad: AdvertisementDetail;
  locale: "fa" | "en";
  t: ReturnType<typeof useI18n>["t"];
}) {
  const media = (ad.media ?? []).filter((item) => item.url || item.thumb);
  const keyValueItems = (entries: [string, unknown][]): DescriptionsProps["items"] =>
    entries.map(([key, value]) => ({
      key,
      label: <span dir="ltr">{key}</span>,
      children: <span dir="auto">{displayValue(value)}</span>,
    }));
  const attrs = keyValueItems(Object.entries(ad.attrs ?? {}));
  const contact = keyValueItems(Object.entries(ad.contact ?? {}));

  const summary: DescriptionsProps["items"] = [
    { key: "status", label: t("status"), children: <AdvertisementStatusTag ad={ad} /> },
    { key: "id", label: t("adId"), children: formatDigits(ad.id, locale) },
    {
      key: "category",
      label: t("category"),
      children: ad.category_title || formatDigits(ad.category_id, locale),
    },
    {
      key: "city",
      label: t("city"),
      children: `${ad.city_name || formatDigits(ad.city_id, locale)}${
        ad.location?.neighborhood ? `، ${ad.location.neighborhood}` : ""
      }`,
    },
    { key: "price", label: t("price"), children: formatPrice(ad, locale, t) },
    {
      key: "customer",
      label: t("customer"),
      children: (
        <>
          {ad.customer_name}{" "}
          <Typography.Text type="secondary" dir="ltr">
            {formatMobile(ad.customer_mobile, locale)}
          </Typography.Text>
        </>
      ),
    },
    {
      key: "created",
      label: t("createdAt"),
      children: formatDateTime(ad.created_at, locale),
    },
    ...(ad.pending_since
      ? [
          {
            key: "pending",
            label: t("waitingSince"),
            children: formatDateTime(ad.pending_since, locale),
          },
        ]
      : []),
  ];

  const reasons: DescriptionsProps["items"] = [
    {
      key: "customer",
      label: t("customerReason"),
      children: <span className="app-pre-wrap">{ad.rejection_reason_customer}</span>,
    },
    {
      key: "internal",
      label: t("internalReason"),
      children: (
        <span className="app-pre-wrap">{ad.rejection_reason_internal ?? "—"}</span>
      ),
    },
  ];

  return (
    <Space orientation="vertical" size="large" style={{ width: "100%" }}>
      {ad.reviewed_at ? (
        <Alert
          type={ad.rejection_reason_customer ? "warning" : "success"}
          showIcon
          title={`${t("lastReview")}: ${ad.reviewed_by_name ?? "—"} · ${formatDateTime(ad.reviewed_at, locale)}`}
          description={
            ad.rejection_reason_customer ? (
              <Descriptions size="small" column={1} items={reasons} />
            ) : null
          }
        />
      ) : null}

      <Descriptions bordered size="small" column={1} items={summary} />

      <section>
        <Typography.Title level={5}>{t("description")}</Typography.Title>
        <Typography.Paragraph className="app-pre-wrap">
          {ad.description}
        </Typography.Paragraph>
      </section>

      <section>
        <Typography.Title level={5}>{t("photos")}</Typography.Title>
        {media.length === 0 ? (
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t("noPhotos")} />
        ) : (
          <Image.PreviewGroup>
            <Space wrap>
              {media.map((item, index) => (
                <Image
                  key={item.url ?? item.thumb ?? index}
                  width={120}
                  height={90}
                  style={{ objectFit: "cover", borderRadius: 6 }}
                  src={mediaUrl(item.thumb ?? item.url)}
                  preview={{ src: mediaUrl(item.url ?? item.thumb) }}
                  alt=""
                />
              ))}
            </Space>
          </Image.PreviewGroup>
        )}
      </section>

      {attrs && attrs.length > 0 ? (
        <section>
          <Typography.Title level={5}>{t("attributes")}</Typography.Title>
          <Descriptions bordered size="small" column={1} items={attrs} />
        </section>
      ) : null}

      {contact && contact.length > 0 ? (
        <section>
          <Typography.Title level={5}>{t("contactInfo")}</Typography.Title>
          <Descriptions bordered size="small" column={1} items={contact} />
        </section>
      ) : null}
    </Space>
  );
}
