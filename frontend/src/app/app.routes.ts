import { Routes } from '@angular/router';
import { UserListComponent } from './features/user-list/user-list.component';
import { UserFormComponent } from './features/user-form/user-form.component';

export const routes: Routes = [
  { path: 'users', component: UserListComponent },
  { path: 'users/new', component: UserFormComponent },
  { path: 'users/edit/:ID', component: UserFormComponent },
  { path: '', redirectTo: '/users', pathMatch: 'full' },
  { path: '**', redirectTo: '/users' }
];
