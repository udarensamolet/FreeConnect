import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { AlertDialogComponent } from '../components/alert-dialog/alert-dialog.component';

@Injectable({
  providedIn: 'root'
})
export class AlertService {
  constructor(private dialog: MatDialog) {}

  alert(message: string): void {
    this.dialog.open(AlertDialogComponent, {
      data: { message: message },
      width: '300px'
    });
  }
}
