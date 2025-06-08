import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class WebsocketService {
  constructor() {}
  private socket: WebSocket | null = null;

  connect(url: string): void {
    if (this.socket) return;

    this.socket = new WebSocket(url);

    this.socket.onopen = () => {
      console.log('WebSocket connected');
    };

    this.socket.onmessage = (event) => {
      console.log('Message from server:', event.data);
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.socket.onclose = () => {
      console.log('WebSocket closed');
      this.socket = null;
    };
  }

  send(message: string) {
    this.socket?.send(message);
  }

  close() {
    this.socket?.close();
  }
}
