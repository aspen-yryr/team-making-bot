import { Injectable } from '@angular/core';
import { FindRequest } from 'proto/match/match_pb';
import { MatchSvcClient } from 'proto/match/MatchServiceClientPb';
import { Match } from '../model/match';

@Injectable({
  providedIn: 'root',
})
export class MatchService {
  client: MatchSvcClient;
  constructor() {
    this.client = new MatchSvcClient('http://localhost:8080');
  }

  async find(id: number): Promise<Match> {
    const req = new FindRequest();
    req.setMatchId(id);
    const match = (await this.client.find(req, null)).toObject().match;
    if (match === undefined) {
      throw new Error('match undefined');
    }
    if (match.owner === undefined) {
      throw new Error('owner undefined');
    }

    return {
      id: match.id,
      owner: {
        id: match.owner.id,
        name: match.owner.name,
      },
      teams: [match.team1, match.team2].map((team) => {
        if (team === undefined) {
          throw new Error('team undefined');
        }
        return {
          id: team.id,
          players: team.playersList.map((p) => {
            return {
              id: p.id,
              name: p.name,
            };
          }),
        };
      }),
    };
  }
}
