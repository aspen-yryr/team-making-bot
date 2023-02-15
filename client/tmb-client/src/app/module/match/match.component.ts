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
    const teams = await this.svc.get(id);
    this.teamA = teams[0];
    this.teamB = teams[1];
    return;
  }
}
