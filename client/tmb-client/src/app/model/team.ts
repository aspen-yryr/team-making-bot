import { User } from './user';

export class Team {
  constructor(readonly id: number, readonly players: User[]) {}
}
