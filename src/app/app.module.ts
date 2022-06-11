import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HeaderComponent } from './header/header.component';
import { HelloComponent } from './hello/hello.component';
import { PicGridComponent } from './pic-grid/pic-grid.component';

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    HelloComponent,
    PicGridComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
