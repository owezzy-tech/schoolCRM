import { Role } from './role';

export class User {
  id!: string;
  img!: string;
  username!: string;
  firstName!: string;
  lastName!: string;
  role!: Role;
  roles!: Role[];
  token!: string;
  expiresAt!: string;
}
