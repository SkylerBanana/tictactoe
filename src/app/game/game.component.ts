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

  winner = signal<string | null>(null);
  board = signal<string[][]>(
    Array(3)
      .fill(null)
      .map(() => Array(3).fill(null))
  );

  startWebSocketX() {
    this.ws.connect('ws://localhost:8085/game?letter=X');
  }
  startWebSocketY() {
    this.ws.connect('ws://localhost:8085/game?letter=Y');
  }

  sendMessage() {
    this.ws.send('Hello from Angular!');
  }

  closeConnection() {
    this.ws.close();
  }
}
