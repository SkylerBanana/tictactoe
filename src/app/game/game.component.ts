import { Component, signal } from '@angular/core';
import { GameBoardComponent } from './game-board/game-board.component';
import { GameInfoComponent } from './game-info/game-info.component';
import { ButtonComponent } from '../components/button/button.component';
import { WebsocketService } from '../websocket.service';

@Component({
  selector: 'app-game',
  standalone: true,
  imports: [GameBoardComponent, GameInfoComponent, ButtonComponent],
  templateUrl: './game.component.html',
  styleUrl: './game.component.css',
})
export class GameComponent {
  constructor(private ws: WebsocketService) {}
  move = signal<Array<Number> | null>(null);
  winner = signal<string | null>(null);
  board = signal<string[][]>(
    Array(3)
      .fill(null)
      .map(() => Array(3).fill(null))
  );

  startWebSocket(letter: string) {
    this.ws.connect(`ws://localhost:8085/game?letter=${letter}`);
  }

  sendMessage() {
    this.ws.send(`${this.move()}`);
  }

  closeConnection() {
    this.ws.close();
  }
  handleClick(i: number, j: number) {
    this.move.set([i, j]);
    this.sendMessage();
  }
}
