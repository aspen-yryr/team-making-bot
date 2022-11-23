import { User } from '../model/user'
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class MatchService {
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

  append(team: string, user: User){
    if (team === "teamA") {
      this.teamA.push(user)
    } else if (team == "teamB"){
      this.teamB.push(user)
    }
    console.log(this.teamA)
    console.log(this.teamB)
  }

  remove(team: string, user: User){
    if (team === "teamA") {
      this.teamA = this.teamA.filter((u) => u.id != user.id)
    } else if (team == "teamB"){
      this.teamB = this.teamB.filter((u) => u.id != user.id)
    }
    console.log(this.teamA)
    console.log(this.teamB)
  }

  get(): User[][] {
    return [this.teamA, this.teamB]
  }

}
