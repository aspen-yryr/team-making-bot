import { User } from '../model/user'
import { Injectable } from '@angular/core';
import { FindRequest } from 'proto/match/match_pb';
import { MatchSvcClient } from 'proto/match/MatchServiceClientPb'

@Injectable({
  providedIn: 'root'
})
export class MatchService {
  client: MatchSvcClient;
  constructor() {
    this.client = new MatchSvcClient('http://localhost:8080')
  }

  async find() {
    const req = new FindRequest();
    req.setMatchId(1);
    return (await this.client.find(req, null)).toObject();
  }
  // TODO: Make user service
  teamA: User[] = [
    new User(1, 'Player1'),
    new User(2, 'Player2'),
  ];
  teamB: User[] = [
    new User(3, 'Player3'),
    new User(4, 'Player4'),
    new User(5, 'Player5'),
  ];

  append(team: string, user: User) {
    if (team === "teamA") {
      this.teamA.push(user)
    } else if (team == "teamB") {
      this.teamB.push(user)
    }
    console.log(this.teamA)
    console.log(this.teamB)
  }

  remove(team: string, user: User) {
    if (team === "teamA") {
      this.teamA = this.teamA.filter((u) => u.id != user.id)
    } else if (team == "teamB") {
      this.teamB = this.teamB.filter((u) => u.id != user.id)
    }
    console.log(this.teamA)
    console.log(this.teamB)
  }

  async get(): Promise<User[][]> {
    const match = await this.find();
    const users = match.match?.team1?.playersList
    if (typeof users == "undefined") {
      return []
    }
    return [users.map<User>((user) => {
      return new User(user.id, user.name)
    }), this.teamB]
  }

}
