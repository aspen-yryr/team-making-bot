import { User } from './user'

export class Team {
    id: number
    players: User[]
    constructor(id: number){
        this.id = id
        this.players = []
    }
}
