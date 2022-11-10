import { writable, Writable } from "svelte/store";
import type { SelfAPI } from "../../../../lib/apiv2";

export interface State {
  messages: object[];
  loading: boolean;
  cursor: number;
}

export class Notifier {
  self_api: SelfAPI;
  state: Writable<State>;

  constructor(self_api: SelfAPI) {
    this.self_api = self_api;
    this.state = writable({ messages: [], cursor: 0, loading: false });
  }

  async init() {
    this.state.update((old) => ({ ...old, loading: true }));

    const resp = await this.self_api.list_message();
    if (resp.status !== 200) {
      console.warn("Error happend", resp);
      return;
    }

    this.state.update((old) => {
      return {
        ...old,
        cursor: 0,
        loading: false,
        messages: resp.data,
      };
    });
  }

  async read_message(id: number) {
    await this.self_api.modify_message({
      ops: "read",
      ids: [id],
    });

    return this.init();
  }
  async delete_message(id: number) {
    await this.self_api.modify_message({
      ops: "delete",
      ids: [id],
    });

    return this.init();
  }
}