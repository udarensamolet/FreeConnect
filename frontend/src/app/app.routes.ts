import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthGuard } from './guards/auth.guard';
import { ClientGuard } from './guards/client.guard';
import { FreelancerGuard } from './guards/freelancer.guard';
import { AdminGuard } from './guards/admin.guard';

import { HomeComponent } from './components/home/home.component';
import { FreelancersComponent } from './components/freelancers/freelancers.component';
import { FreelancerProfileComponent } from './components/freelancer-profile/freelancer-profile.component';
import { RegisterComponent } from './components/register/register.component';
import { LoginComponent } from './components/login/login.component';
import { PostProjectComponent } from './components/post-project/post-project.component';
import { MyProjectsComponent } from './components/my-projects/my-projects.component';
import { ProposalsComponent } from './components/proposals/proposals.component';
import { TaskDetailsComponent } from './components/task-details/task-details.component';
import { ClientDashboardComponent } from './components/client-dashboard/client-dashboard.component';
import { FreelancerDashboardComponent } from './components/freelancer-dashboard/freelancer-dashboard.component';
import { FreelancerProjectsComponent } from './components/freelancer-projects/freelancer-projects.component';
import { AdminDashboardComponent } from './components/admin-dashboard/admin-dashboard.component';
import { EditProjectComponent } from './components/edit-project/edit-project.component';
import { ProjectDetailsComponent } from './components/project-details/project-details.component';
import { EditProfileComponent } from './components/edit-profile/edit-profile.component';
import { EditTaskComponent } from './components/edit-task/edit-task.component';
import { CreateTaskComponent } from './components/create-task/create-task.component';
import { ProfileDetailsComponent } from './components/profile-details/profile-details.component';
import { BrowseUnassignedProjectsComponent } from './components/browse-projects/browse-projects.component';

export const routes: Routes = [
  // Default route -> go to /home
  { path: '', redirectTo: 'home', pathMatch: 'full' },

  // Public or basic routes
  { path: 'home', component: HomeComponent },
  { path: 'freelancers', component: FreelancersComponent },
  { path: 'freelancers/:id', component: FreelancerProfileComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },

  // Client routes
  {
    path: 'post-project',
    component: PostProjectComponent,
    canActivate: [ClientGuard]
  },
  {
    path: 'my-projects',
    component: MyProjectsComponent,
    canActivate: [ClientGuard]
  },
  {
    path: 'projects/:id/edit',
    component: EditProjectComponent,
    canActivate: [ClientGuard]
  },
  {
    path: 'projects/:id/create-task',
    component: CreateTaskComponent,
    canActivate: [ClientGuard]
  },
  {
    path: 'projects/:projectId/tasks/:taskId/edit',
    component: EditTaskComponent,
    canActivate: [ClientGuard]
  },
  {
    path: 'projects/:projectId/tasks/:taskId',
    component: TaskDetailsComponent,
    canActivate: [ClientGuard]
  },

  // Shared or open routes
  {
    path: 'projects/:id',
    component: ProjectDetailsComponent
  },
  {
    path: 'proposals',
    component: ProposalsComponent
  },

  // Freelancer routes
  {
    path: 'freelancer-projects',
    component: FreelancerProjectsComponent,
    canActivate: [FreelancerGuard]
  },

  // Admin
  {
    path: 'admin-dashboard',
    component: AdminDashboardComponent,
    canActivate: [AdminGuard]
  },

  // Dashboard examples with roles
  {
    path: 'client-dashboard',
    component: ClientDashboardComponent,
    canActivate: [AuthGuard],
    data: { roles: ['client'] }
  },
  {
    path: 'freelancer-dashboard',
    component: FreelancerDashboardComponent,
    canActivate: [AuthGuard],
    data: { roles: ['freelancer'] }
  },

  // Profile
  {
    path: 'profile',
    component: ProfileDetailsComponent,
    canActivate: [AuthGuard]
  },
  {
    path: 'profile/edit',
    component: EditProfileComponent,
    canActivate: [AuthGuard]
  },

  { path: '**', redirectTo: 'home', pathMatch: 'full' },
  {
    path: 'browse-projects',
    component: BrowseUnassignedProjectsComponent
  }
  
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}