import { Component, inject } from '@angular/core';
import { FormsModule, NgForm } from '@angular/forms';
import { HttpClient, HttpParams } from '@angular/common/http';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
})
export class LoginComponent {
  isLoggedin = true;
  username = '';
  private http = inject(HttpClient);

  // There is Probrably a better way of doing this

  // what this essentially does is it checks if we have a cookie and since we cant check http only cookies with javascript it calls the server to check

  ngOnInit() {
    this.http
      .get('http://localhost:8085/checkuser', {
        withCredentials: true, // ill delete this in prod
      })
      .subscribe({
        next: () => {
          this.isLoggedin = true;
        },
        error: () => {
          this.isLoggedin = false;
        },
      });
  }

  onSubmit(form: NgForm) {
    if (form.valid) {
      console.log(this.username);
      const body = new HttpParams().set('username', this.username);

      this.http
        .post('http://localhost:8085/login', body.toString(), {
          headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          withCredentials: true, // ill delete this in prod
        })
        .subscribe({
          next: (res) => {
            console.log('Login Successful', res);
            this.isLoggedin = true;
          },
          error: (err) => {
            console.error('Login Failed', err);
            this.isLoggedin = false;
          },
        });
    }
  }
}
