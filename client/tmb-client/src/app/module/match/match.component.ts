import { Component, OnInit } from '@angular/core';
import { MatchService } from 'src/app/service/match.service';
import { ActivatedRoute } from '@angular/router';
import { User } from 'src/app/model/user';

@Component({
  selector: 'app-match',
  templateUrl: './match.component.html',
  styleUrls: ['./match.component.css'],
})
export class MatchComponent implements OnInit {
  teamA: User[] = [];
  teamB: User[] = [];

  constructor(
    private readonly svc: MatchService,
    private readonly route: ActivatedRoute,
  ) {}

  async ngOnInit(): Promise<void> {
    const id = Number(this.route.snapshot.paramMap.get('id'));
    const match = await this.svc.find(id);
    this.teamA = match.teams[0].players;
    this.teamB = match.teams[1].players;
    return;
  }
}
