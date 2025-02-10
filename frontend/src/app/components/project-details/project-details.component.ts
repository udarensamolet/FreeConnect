// import { Component, OnInit } from '@angular/core';
// import { ActivatedRoute } from '@angular/router';
// import { ProjectService } from '../../services/project.service';
// import { AuthService } from '../../services/auth.service';
// import { ProposalService } from '../../services/proposal.service';

// @Component({
//   selector: 'app-project-details',
//   templateUrl: './project-details.component.html',
//   styleUrls: ['./project-details.component.css']
// })
// export class ProjectDetailsComponent implements OnInit {
//   user: any;
//   proposalText = '';
//   estimatedDuration = 0;
//   bidAmount = 0;

//   constructor(
//     private route: ActivatedRoute,
//     private projectService: ProjectService,
//     private proposalService: ProposalService,
//     private authService: AuthService
//   ) {}

//   ngOnInit(): void {
//     const projectId = Number(this.route.snapshot.paramMap.get('id'));
//     this.loadProject(projectId);
//   }

//   loadProject(id: number): void {
//     this.projectService.getProjectById(id).subscribe((response: any) => {
//       this.project = response.project;
//     });
//   }
// }