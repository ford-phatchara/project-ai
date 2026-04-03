import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { UserService } from '../../core/user.service';
import { User } from '../../models/user.model';
import { finalize } from 'rxjs';

@Component({
  selector: 'app-user-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  template: `
    <div class="container">
      <h2>{{ isEditMode ? 'Update User' : 'Create New User' }}</h2>

      <div *ngIf="loading" class="loading">Loading...</div>
      <div *ngIf="error" class="error">{{ error }}</div>

      <form [formGroup]="userForm" (ngSubmit)="onSubmit()" *ngIf="!loading">
        <div class="form-group">
          <label for="name">Name</label>
          <input id="name" type="text" formControlName="name" class="form-control" placeholder="Enter name">
          <div *ngIf="f['name'].touched && f['name'].invalid" class="invalid-feedback">
            Name is required (min 2 characters).
          </div>
        </div>

        <div class="form-group">
          <label for="email">Email</label>
          <input id="email" type="email" formControlName="email" class="form-control" placeholder="Enter email">
          <div *ngIf="f['email'].touched && f['email'].invalid" class="invalid-feedback">
            Please enter a valid email.
          </div>
        </div>

        <div class="btn-group">
          <button type="submit" [disabled]="userForm.invalid || submitting" class="btn btn-primary">
            {{ submitting ? 'Submitting...' : (isEditMode ? 'Update' : 'Create') }}
          </button>
          <button type="button" routerLink="/users" class="btn btn-secondary">Cancel</button>
        </div>
      </form>
    </div>
  `,
  styles: [`
    .container { padding: 20px; max-width: 500px; }
    .form-group { margin-bottom: 15px; }
    .form-group label { display: block; margin-bottom: 5px; }
    .form-control { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
    .btn-group { margin-top: 20px; }
    .btn { padding: 8px 16px; cursor: pointer; border-radius: 4px; border: none; }
    .btn-primary { background-color: #007bff; color: white; margin-right: 10px; }
    .btn-primary:disabled { background-color: #ccc; cursor: not-allowed; }
    .btn-secondary { background-color: #6c757d; color: white; }
    .invalid-feedback { color: red; font-size: 0.8em; margin-top: 5px; }
    .loading { margin-top: 10px; font-style: italic; }
    .error { color: red; margin-top: 10px; }
  `]
})
export class UserFormComponent implements OnInit {
  private fb = inject(FormBuilder);
  private userService = inject(UserService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  userForm: FormGroup = this.fb.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    email: ['', [Validators.required, Validators.email]]
  });

  userId: number | null = null;
  isEditMode = false;
  loading = false;
  submitting = false;
  error: string | null = null;

  get f() { return this.userForm.controls; }

  ngOnInit(): void {
    // Use paramMap to handle parameter changes if component is reused
    this.route.paramMap.subscribe(params => {
      const ID = params.get('ID');
      if (ID) {
        this.userId = +ID;
        this.isEditMode = true;
        this.loadUser(this.userId);
      } else {
        this.userId = null;
        this.isEditMode = false;
        this.userForm.reset();
        this.error = null;
      }
    });
  }

  loadUser(id: number): void {
    // this.loading = true;
    this.userService.getUserById(id).pipe(
      finalize(() => this.loading = false)
    ).subscribe({
      next: (user) => {
        console.log('Loaded user: ', user);
        this.userForm.patchValue(user);
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load user data.';
        this.loading = false;
      }
    });
  }

  onSubmit(): void {
    if (this.userForm.invalid) return;

    this.submitting = true;
    this.error = null;

    const userData: User = this.userForm.value;

    if (this.isEditMode && this.userId) {
      this.userService.updateUser(this.userId, userData).pipe(
        finalize(() => this.submitting = false)
      ).subscribe({
        next: () => {
          this.router.navigate(['/users']);
        },
        error: (err) => {
          this.error = 'Failed to update user.';
        }
      });
    } else {
      this.userService.createUser(userData).pipe(
        finalize(() => this.submitting = false)
      ).subscribe({
        next: () => {
          this.router.navigate(['/users']);
        },
        error: (err) => {
          this.error = 'Failed to create user.';
        }
      });
    }
  }
}
