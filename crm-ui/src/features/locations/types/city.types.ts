export interface City {
  id: number;
  name: string;
  /** Province name, to tell apart cities that share a name. */
  province: string | null;
}
