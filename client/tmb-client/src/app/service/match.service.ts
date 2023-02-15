import { User } from '../model/user';
import { Injectable } from '@angular/core';
import { FindRequest } from 'proto/match/match_pb';
import { MatchSvcClient } from 'proto/match/MatchServiceClientPb';

@Injectable({
  providedIn: 'root',
})
export class MatchService {
  client: MatchSvcClient;
  constructor() {
    this.client = new MatchSvcClient('http://localhost:8080');
  }

  async find(id: number) {
    const req = new FindRequest();
    req.setMatchId(id);
    return (await this.client.find(req, null)).toObject();
  }

  async get(id: number): Promise<User[][]> {
    const match = await this.find(id);

    const users1 = match.match?.team1?.playersList;
    if (typeof users1 == 'undefined') {
      return [];
    }
    const users2 = match.match?.team2?.playersList;
    if (typeof users2 == 'undefined') {
      return [];
    }

    return [
      users1.map<User>((user) => {
        return new User(user.id, user.name);
      }),
      users2.map<User>((user) => {
        return new User(user.id, user.name);
      }),
    ];
  }
}
