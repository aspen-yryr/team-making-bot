import * as jspb from 'google-protobuf'

import * as proto_match_user_pb from '../../proto/match/user_pb';


export class Team extends jspb.Message {
  getId(): number;
  setId(value: number): Team;

  getPlayersList(): Array<proto_match_user_pb.User>;
  setPlayersList(value: Array<proto_match_user_pb.User>): Team;
  clearPlayersList(): Team;
  addPlayers(value?: proto_match_user_pb.User, index?: number): proto_match_user_pb.User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Team.AsObject;
  static toObject(includeInstance: boolean, msg: Team): Team.AsObject;
  static serializeBinaryToWriter(message: Team, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Team;
  static deserializeBinaryFromReader(message: Team, reader: jspb.BinaryReader): Team;
}

export namespace Team {
  export type AsObject = {
    id: number,
    playersList: Array<proto_match_user_pb.User.AsObject>,
  }
}

export class Match extends jspb.Message {
  getId(): number;
  setId(value: number): Match;

  getOwner(): proto_match_user_pb.User | undefined;
  setOwner(value?: proto_match_user_pb.User): Match;
  hasOwner(): boolean;
  clearOwner(): Match;

  getMembersList(): Array<proto_match_user_pb.User>;
  setMembersList(value: Array<proto_match_user_pb.User>): Match;
  clearMembersList(): Match;
  addMembers(value?: proto_match_user_pb.User, index?: number): proto_match_user_pb.User;

  getTeam1(): Team | undefined;
  setTeam1(value?: Team): Match;
  hasTeam1(): boolean;
  clearTeam1(): Match;

  getTeam2(): Team | undefined;
  setTeam2(value?: Team): Match;
  hasTeam2(): boolean;
  clearTeam2(): Match;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Match.AsObject;
  static toObject(includeInstance: boolean, msg: Match): Match.AsObject;
  static serializeBinaryToWriter(message: Match, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Match;
  static deserializeBinaryFromReader(message: Match, reader: jspb.BinaryReader): Match;
}

export namespace Match {
  export type AsObject = {
    id: number,
    owner?: proto_match_user_pb.User.AsObject,
    membersList: Array<proto_match_user_pb.User.AsObject>,
    team1?: Team.AsObject,
    team2?: Team.AsObject,
  }
}

export class CreateUserRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserRequest): CreateUserRequest.AsObject;
  static serializeBinaryToWriter(message: CreateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserRequest;
  static deserializeBinaryFromReader(message: CreateUserRequest, reader: jspb.BinaryReader): CreateUserRequest;
}

export namespace CreateUserRequest {
  export type AsObject = {
    name: string,
  }
}

export class CreateUserResponse extends jspb.Message {
  getUser(): proto_match_user_pb.User | undefined;
  setUser(value?: proto_match_user_pb.User): CreateUserResponse;
  hasUser(): boolean;
  clearUser(): CreateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserResponse): CreateUserResponse.AsObject;
  static serializeBinaryToWriter(message: CreateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserResponse;
  static deserializeBinaryFromReader(message: CreateUserResponse, reader: jspb.BinaryReader): CreateUserResponse;
}

export namespace CreateUserResponse {
  export type AsObject = {
    user?: proto_match_user_pb.User.AsObject,
  }
}

export class CreateMatchRequest extends jspb.Message {
  getOwner(): proto_match_user_pb.User | undefined;
  setOwner(value?: proto_match_user_pb.User): CreateMatchRequest;
  hasOwner(): boolean;
  clearOwner(): CreateMatchRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateMatchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateMatchRequest): CreateMatchRequest.AsObject;
  static serializeBinaryToWriter(message: CreateMatchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateMatchRequest;
  static deserializeBinaryFromReader(message: CreateMatchRequest, reader: jspb.BinaryReader): CreateMatchRequest;
}

export namespace CreateMatchRequest {
  export type AsObject = {
    owner?: proto_match_user_pb.User.AsObject,
  }
}

export class CreateMatchResponse extends jspb.Message {
  getMatch(): Match | undefined;
  setMatch(value?: Match): CreateMatchResponse;
  hasMatch(): boolean;
  clearMatch(): CreateMatchResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateMatchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateMatchResponse): CreateMatchResponse.AsObject;
  static serializeBinaryToWriter(message: CreateMatchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateMatchResponse;
  static deserializeBinaryFromReader(message: CreateMatchResponse, reader: jspb.BinaryReader): CreateMatchResponse;
}

export namespace CreateMatchResponse {
  export type AsObject = {
    match?: Match.AsObject,
  }
}

export class FindRequest extends jspb.Message {
  getMatchId(): number;
  setMatchId(value: number): FindRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FindRequest): FindRequest.AsObject;
  static serializeBinaryToWriter(message: FindRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindRequest;
  static deserializeBinaryFromReader(message: FindRequest, reader: jspb.BinaryReader): FindRequest;
}

export namespace FindRequest {
  export type AsObject = {
    matchId: number,
  }
}

export class FindResponse extends jspb.Message {
  getMatch(): Match | undefined;
  setMatch(value?: Match): FindResponse;
  hasMatch(): boolean;
  clearMatch(): FindResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FindResponse): FindResponse.AsObject;
  static serializeBinaryToWriter(message: FindResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindResponse;
  static deserializeBinaryFromReader(message: FindResponse, reader: jspb.BinaryReader): FindResponse;
}

export namespace FindResponse {
  export type AsObject = {
    match?: Match.AsObject,
  }
}

export class AppendMemberRequest extends jspb.Message {
  getMatchId(): number;
  setMatchId(value: number): AppendMemberRequest;

  getMembersList(): Array<proto_match_user_pb.User>;
  setMembersList(value: Array<proto_match_user_pb.User>): AppendMemberRequest;
  clearMembersList(): AppendMemberRequest;
  addMembers(value?: proto_match_user_pb.User, index?: number): proto_match_user_pb.User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AppendMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AppendMemberRequest): AppendMemberRequest.AsObject;
  static serializeBinaryToWriter(message: AppendMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AppendMemberRequest;
  static deserializeBinaryFromReader(message: AppendMemberRequest, reader: jspb.BinaryReader): AppendMemberRequest;
}

export namespace AppendMemberRequest {
  export type AsObject = {
    matchId: number,
    membersList: Array<proto_match_user_pb.User.AsObject>,
  }
}

export class ShuffleRequest extends jspb.Message {
  getMatchId(): number;
  setMatchId(value: number): ShuffleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShuffleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ShuffleRequest): ShuffleRequest.AsObject;
  static serializeBinaryToWriter(message: ShuffleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShuffleRequest;
  static deserializeBinaryFromReader(message: ShuffleRequest, reader: jspb.BinaryReader): ShuffleRequest;
}

export namespace ShuffleRequest {
  export type AsObject = {
    matchId: number,
  }
}

export class ShuffleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShuffleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ShuffleResponse): ShuffleResponse.AsObject;
  static serializeBinaryToWriter(message: ShuffleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShuffleResponse;
  static deserializeBinaryFromReader(message: ShuffleResponse, reader: jspb.BinaryReader): ShuffleResponse;
}

export namespace ShuffleResponse {
  export type AsObject = {
  }
}

