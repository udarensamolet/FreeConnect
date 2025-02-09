import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from './guards/auth.guard';

import { HomeComponent } from './components/home/home.component';
import { ProjectsComponent } from './components/projects/projects.component';
import { ProjectDetailsComponent } from './components/project-details/project-details.component';
import { FreelancersComponent } from './components/freelancers/freelancers.component';
import { FreelancerProfileComponent } from './components/freelancer-profile/freelancer-profile.component';
import { RegisterComponent } from './components/register/register.component';
import { LoginComponent } from './components/login/login.component';
import { UserProfileComponent } from './components/user-profile/user-profile.component';
import { PostProjectComponent } from './components/post-project/post-project.component';
import { MyProjectsComponent } from './components/my-projects/my-projects.component';
import { ProposalsComponent } from './components/proposals/proposals.component';
import { TaskDetailsComponent } from './components/task-details/task-details.component';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },

  { path: 'home', component: HomeComponent, canActivate: [AuthGuard] },
  { path: 'projects', component: ProjectsComponent },
  { path: 'projects/:id', component: ProjectDetailsComponent },
  { path: 'freelancers', component: FreelancersComponent },
  { path: 'freelancers/:id', component: FreelancerProfileComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },
  { path: 'profile', component: UserProfileComponent },
  { path: 'post-project', component: PostProjectComponent },
  { path: 'my-projects', component: MyProjectsComponent },
  { path: 'proposals', component: ProposalsComponent },
  { path: 'projects/:id/tasks/:taskId', component: TaskDetailsComponent },

  { path: '**', redirectTo: 'home', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
