import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';

import { UserService } from '../../services/user.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-add-skill-dialog',
  standalone: true,
  templateUrl: './add-skill-dialog.component.html',
  styleUrls: ['./add-skill-dialog.component.css'],
  imports: [
    CommonModule,
    FormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ]
})
export class AddSkillDialogComponent {
  skillName = '';

  constructor(
    private dialogRef: MatDialogRef<AddSkillDialogComponent>,
    private userService: UserService,
    private authService: AuthService
  ) {}

  onAddSkill() {
    if (!this.skillName.trim()) {
      return;
    }
    // 1) get the current user
    const user = this.authService.getCurrentUser();
    if (!user) {
      this.dialogRef.close();
      return;
    }

    // 2) call userService to add skill
    this.userService.addSkillToUser(user.user_id, this.skillName).subscribe({
      next: () => {
        // close with 'skillAdded' so the parent can reload
        this.dialogRef.close('skillAdded');
      },
      error: (err) => {
        console.error('Failed to add skill:', err);
      }
    });
  }

  onCancel() {
    this.dialogRef.close();
  }
}