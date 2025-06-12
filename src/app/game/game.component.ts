import { Component, signal, effect } from '@angular/core';
import { GameBoardComponent } from './game-board/game-board.component';
import { GameInfoComponent } from './game-info/game-info.component';
import { ButtonComponent } from '../components/button/button.component';
import { WebsocketService } from '../websocket.service';
import { LoginComponent } from '../components/login/login.component';

@Component({
  selector: 'app-game',
  standalone: true,
  imports: [
    GameBoardComponent,
    GameInfoComponent,
    ButtonComponent,
    LoginComponent,
  ],
  templateUrl: './game.component.html',
  styleUrl: './game.component.css',
})
export class GameComponent {
  constructor(private ws: WebsocketService) {
    effect(() => {
      // Pretty much if server message changes This code runs
      const message = this.ws.serverMessage();
      if (message?.matrix) {
        this.board.set(message.matrix);
        console.log(this.board);
      }
    });
  }
  move = signal<Array<Number> | null>(null);
  winner = signal<string | null>(null);
  board = signal<string[][]>(
    Array(3)
      .fill(null)
      .map(() => Array(3).fill(null))
  );

  startWebSocket() {
    this.ws.connect(`ws://localhost:8085/game`);
  }

  sendMessage() {
    this.ws.send(JSON.stringify(this.move()));
  }

  closeConnection() {
    this.ws.close();
  }
  handleClick(i: number, j: number) {
    this.move.set([i, j]);
    this.sendMessage();
  }
}
