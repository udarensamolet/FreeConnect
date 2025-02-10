import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';

import { AuthService } from '../../services/auth.service';
import { UserService } from '../../services/user.service';
import { AddSkillDialogComponent } from '../add-skill-dialog/add-skill-dialog.component';

@Component({
  selector: 'app-profile-details',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule
  ],
  templateUrl: './profile-details.component.html',
  styleUrls: ['./profile-details.component.css']
})
export class ProfileDetailsComponent implements OnInit {
  user: any = null;

  constructor(
    private authService: AuthService,
    private userService: UserService,
    private router: Router,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) {
      this.router.navigate(['/login']);
      return;
    }
    this.loadUser(currentUser.user_id);
  }

  loadUser(userId: number): void {
    this.userService.getUserById(userId).subscribe({
      next: (res: any) => {
        this.user = res.user || res;
      },
      error: (err) => {
        console.error('Failed to load user:', err);
      }
    });
  }

  onEditProfile(): void {
    this.router.navigate(['/profile', 'edit']);
  }

  onAddSkill(): void {
    const dialogRef = this.dialog.open(AddSkillDialogComponent, {
      width: '350px'
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result === 'skillAdded') {
        this.reloadProfile();
      }
    });
  }

  reloadProfile(): void {
    if (this.user) {
      this.loadUser(this.user.user_id);
    }
  }
}