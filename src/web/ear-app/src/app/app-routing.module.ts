import { ResultsComponent } from './results/results.component';
import { UploadComponent } from './upload/upload.component';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

const routes: Routes = [
  {path: 'upload', component: UploadComponent},
  {path: 'results', component: ResultsComponent},
  {path: 'results/:id', component: ResultsComponent},
  {
    path: '',
    redirectTo: '/upload',
    pathMatch: 'full'
  },
  {
    path: '**',
    redirectTo: '/upload',
    pathMatch: 'full'
  },

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
