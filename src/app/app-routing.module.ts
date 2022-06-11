import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HelloComponent } from './hello/hello.component';
import { PicGridComponent } from './pic-grid/pic-grid.component';

const routes: Routes = [{ path: 'home', component: HelloComponent }, { path: '**', redirectTo: '/home' }];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
