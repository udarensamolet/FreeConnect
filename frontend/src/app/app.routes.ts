import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthGuard } from './guards/auth.guard';
import { ClientGuard } from './guards/client.guard';
import { FreelancerGuard } from './guards/freelancer.guard';
import { AdminGuard } from './guards/admin.guard';

import { HomeComponent } from './components/home/home.component';
// import { ProjectsComponent } from './components/projects/projects.component';
// import { ProjectDetailsComponent } from './components/project-details/project-details.component';
import { FreelancersComponent } from './components/freelancers/freelancers.component';
import { FreelancerProfileComponent } from './components/freelancer-profile/freelancer-profile.component';
import { RegisterComponent } from './components/register/register.component';
import { LoginComponent } from './components/login/login.component';
import { UserProfileComponent } from './components/user-profile/user-profile.component';
import { PostProjectComponent } from './components/post-project/post-project.component';
import { MyProjectsComponent } from './components/my-projects/my-projects.component';
import { ProposalsComponent } from './components/proposals/proposals.component';
import { TaskDetailsComponent } from './components/task-details/task-details.component';
import { ClientDashboardComponent } from './components/client-dashboard/client-dashboard.component';
import { FreelancerDashboardComponent } from './components/freelancer-dashboard/freelancer-dashboard.component';
import { FreelancerProjectsComponent } from './components/freelancer-projects/freelancer-projects.component';
import { AdminDashboardComponent } from './components/admin-dashboard/admin-dashboard.component';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: '', component: HomeComponent },
  { path: 'home', component: HomeComponent },
  // { 
  //   path: 'projects', 
  //   component: ProjectsComponent, 
  //   canActivate: [AuthGuard], 
  //   data: { roles: ['client', 'freelancer', 'admin']}
  // },
  // { path: 'projects/:id', component: ProjectDetailsComponent },
  { path: 'freelancers', component: FreelancersComponent },
  { path: 'freelancers/:id', component: FreelancerProfileComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },
  { path: 'profile', component: UserProfileComponent },
  { path: 'post-project', component: PostProjectComponent },
  {
    path: 'my-projects',
    component: MyProjectsComponent,
    canActivate: [ClientGuard]
  },
  {
    path: 'freelancer-projects',
    component: FreelancerProjectsComponent,
    canActivate: [FreelancerGuard]
  },
  {
    path: 'admin-dashboard',
    component: AdminDashboardComponent,
    canActivate: [AdminGuard]
  },
  // { path: 'project/:id', component: ProjectDetailsComponent },
  { path: 'proposals', component: ProposalsComponent },
  { path: 'projects/:id/tasks/:taskId', component: TaskDetailsComponent },
  { path: '**', redirectTo: 'home', pathMatch: 'full' },
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
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
