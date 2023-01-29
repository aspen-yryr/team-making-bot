import {
  CdkDragDrop,
  moveItemInArray,
} from '@angular/cdk/drag-drop';
import { Component, OnInit } from '@angular/core';
import { Team } from '../../model/team'
import { MatchService } from 'src/app/service/match.service';

@Component({
  selector: 'app-match',
  templateUrl: './match.component.html',
  styleUrls: ['./match.component.css']
})
export class MatchComponent implements OnInit {
  teamA = new Team(0);
  teamB = new Team(1);

  constructor(private readonly svc: MatchService) { }

  // TODO: Use ngrx
  async drop(event: CdkDragDrop<Team>) {
    if (event.previousContainer === event.container) {
      moveItemInArray(
        event.container.data.players,
        event.previousIndex,
        event.currentIndex
      );
    } else {
      if (event.previousContainer.data.id === 0) {
        this.svc.remove('teamA', event.previousContainer.data.players[event.previousIndex])
        this.svc.append('teamB', event.previousContainer.data.players[event.previousIndex])
      } else if (event.previousContainer.data.id === 1) {
        this.svc.remove('teamB', event.previousContainer.data.players[event.previousIndex])
        this.svc.append('teamA', event.previousContainer.data.players[event.previousIndex])
      }
      const teams = await this.svc.get()
      this.teamA.players = teams[0]
      this.teamB.players = teams[1]
    }
  }

  async ngOnInit(): Promise<void> {
    const teams = await this.svc.get()
    this.teamA.players = teams[0]
    this.teamB.players = teams[1]
    return
  }
}
