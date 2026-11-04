export interface Category {
  id: number;
  parent: number | null;
  order: number;
  title: string;
  slug: string;
  isLeaf: boolean;
}
