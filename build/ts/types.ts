export interface ICard {
  readonly id: number;
  readonly assignees: IMembers;
  readonly in_play: number;
  readonly status: string;
  readonly stickers: any;
  readonly title: string;
  readonly url: string;
}

export interface IMembers {
  readonly [id: string]: IMember;
}

interface IMember {
  readonly id: number;
  readonly email: string;
  readonly name: string;
}

export interface ISupport {
  readonly type: string; //tslint:disable-line:no-reserved-keywords
  readonly member: string;
}
