import { User } from './user';
import { Team } from './team';

export class Match {
  constructor(
    readonly id: number,
    readonly owner: User,
    readonly teams: Team[]
  ) {}
}
