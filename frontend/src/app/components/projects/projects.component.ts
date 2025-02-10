// // projects.component.ts
// import { Component, OnInit } from '@angular/core';
// import { ProjectService } from '../../services/project.service';

// @Component({
//   selector: 'app-projects',
//   templateUrl: './projects.component.html',
//   styleUrls: ['./projects.component.css']
// })
// export class ProjectsComponent implements OnInit {
//   projects: any[] = [];

//   searchText = '';
//   minBudget?: number;
//   maxBudget?: number;
//   status?: string;

//   constructor(private projectService: ProjectService) {}

//   ngOnInit(): void {
//     this.loadProjects();
//   }

//   loadProjects(): void {
//     this.projectService
//       .getProjects(this.searchText, this.minBudget, this.maxBudget, this.status)
//       .subscribe((data: any) => {
//         this.projects = data.projects;
//       });
//   }

//   onSearch(): void {
//     this.loadProjects();
//   }
// }
