import $ from 'jquery';
import { ICard, IMembers, ISupport } from './types';

interface IResponse {
  readonly cards: ReadonlyArray<ICard>;
  readonly support: {
    readonly [type: string]: ISupport; // tslint:disable-line:no-reserved-keywords
  };
  readonly free_team_members: IMembers;
}

export default class State {
  public static updated: string = 'rubbernecker:state:updated';

  public content: IResponse;

  constructor() {
    this.content = {
      cards: [],
      support: {},
      free_team_members: {},
    };
  }

  public async fetchState() {
    const request = new Request('/state', {
      method: 'GET',
      headers: new Headers({
        Accept: 'application/json',
      }),
    });

    const response = await fetch(request);

    if (response.status !== 200) {
      console.error('Rubbernecker responded with non 200 http status.');
    }

    const data: any = await response.json();

    this.content = data;

    // Trigger updated event.
    $(document).trigger(State.updated);
  }
}
