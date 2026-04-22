import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login-page',
  standalone: true,
  imports: [ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <main class="flex min-h-screen items-center justify-center px-4 py-10">
      <section class="w-full max-w-md rounded-xl border border-stone-200 bg-white/95 p-6 shadow-xl shadow-stone-900/10 sm:p-8">
        <div class="mb-8">
          <div class="mb-4 flex size-12 items-center justify-center rounded-xl bg-emerald-700 text-white shadow-sm">
            <svg class="size-6" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M12 3.5c4.25 2.1 6.5 5.03 6.5 8.47 0 4.57-3.4 8.53-6.5 8.53s-6.5-3.96-6.5-8.53C5.5 8.53 7.75 5.6 12 3.5Z"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linejoin="round"
              />
              <path d="M12 7.5v9" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
            </svg>
          </div>
          <p class="text-sm font-semibold uppercase tracking-wide text-emerald-700">
            Durian Farm Manager
          </p>
          <h1 class="mt-2 text-3xl font-semibold tracking-tight text-stone-950">
            Welcome back
          </h1>
          <p class="mt-3 text-sm leading-6 text-stone-500">
            Sign in to manage sales, plots, and care activity.
          </p>
        </div>

        <form class="space-y-5" [formGroup]="form" (ngSubmit)="login()">
          <label class="block space-y-2">
            <span class="form-label">Email</span>
            <input
              class="form-field"
              type="email"
              autocomplete="email"
              formControlName="email"
              placeholder="manager@farm.test"
            />
            @if (form.controls.email.invalid && form.controls.email.touched) {
              <span class="text-sm text-red-600">Enter a valid email.</span>
            }
          </label>

          <label class="block space-y-2">
            <span class="form-label">Password</span>
            <input
              class="form-field"
              type="password"
              autocomplete="current-password"
              formControlName="password"
              placeholder="password"
            />
            @if (form.controls.password.invalid && form.controls.password.touched) {
              <span class="text-sm text-red-600">Password is required.</span>
            }
          </label>

          <button class="btn-primary w-full" type="submit">Login</button>
        </form>
      </section>
    </main>
  `,
})
export class LoginPageComponent {
  private readonly fb = inject(NonNullableFormBuilder);
  private readonly router = inject(Router);

  readonly form = this.fb.group({
    email: ['manager@durianfarm.test', [Validators.required, Validators.email]],
    password: ['durian2026', Validators.required],
  });

  login(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    void this.router.navigateByUrl('/dashboard');
  }
}
