export interface NamedEntity {
  id: number;
  code: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export type IamAction = NamedEntity;
export type IamRole = NamedEntity;
export type IamGroup = NamedEntity;

export interface NamedEntityInput {
  code: string;
  name: string;
}
