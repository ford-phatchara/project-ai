import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { Observable, finalize } from 'rxjs';
import { UserService } from '../../core/user.service';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="container">
      <h2>User List</h2>
      <button routerLink="/users/new" class="btn btn-primary">Add New User</button>
      <div *ngIf="loading" class="loading">Loading users...</div>
      <div *ngIf="error" class="error">{{ error }}</div>

      <table *ngIf="!loading && !error" class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr *ngFor="let user of users$ | async">
            <td>{{ user.ID }}</td>
            <td>{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>
              <button [routerLink]="['/users/edit',user.ID]" class="btn btn-secondary">Edit</button>
            </td>
          </tr>
        </tbody>
      </table>

      <div *ngIf="(users$ | async)?.length === 0 && !loading" class="no-users">
        No users found.
      </div>
    </div>
  `,
  styles: [`
    .container { padding: 20px; }
    .table { width: 100%; border-collapse: collapse; margin-top: 20px; }
    .table th, .table td { border: 1px solid #ddd; padding: 8px; text-align: left; }
    .table th { background-color: #f2f2f2; }
    .btn { padding: 8px 16px; cursor: pointer; border-radius: 4px; border: none; }
    .btn-primary { background-color: #007bff; color: white; }
    .btn-secondary { background-color: #6c757d; color: white; margin-left: 5px; }
    .loading { margin-top: 10px; font-style: italic; }
    .error { color: red; margin-top: 10px; }
    .no-users { margin-top: 10px; }
  `]
})
export class UserListComponent implements OnInit {
  private userService = inject(UserService);
  users$!: Observable<User[]>;
  loading = true;
  error: string | null = null;

  ngOnInit(): void {
    this.loadUsers();
  }

  loadUsers(): void {
    this.loading = true;
    this.users$ = this.userService.getUsers()
      .pipe(
        finalize(() => this.loading = false)
      );
      console.log('this.users$ ',this.users$);
  }
}
