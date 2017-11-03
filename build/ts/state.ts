import { ICard, IMembers, ISupport } from './types';

interface IResponse {
  readonly cards: ReadonlyArray<ICard>;
  readonly support: {
    readonly [type: string]: ISupport; // tslint:disable-line:no-reserved-keywords
  };
  readonly free_team_members: IMembers;
}

declare var $: any;

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

  public fetchState() {
    const request = new Request('/state', {
      method: 'GET',
      headers: new Headers({
        Accept: 'application/json',
      }),
    });

    fetch(request)
      .then(response => {
        if (response.status !== 200) {
          console.error('Rubbernecker responded with non 200 http status.');
        }

        return response.json();
      })
      .then(data => {
        this.content = data;

        // Trigger updated event.
        $(document).trigger(State.updated);
      });
  }
}
