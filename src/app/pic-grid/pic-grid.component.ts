import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';

@Component({
  selector: 'app-pic-grid',
  templateUrl: './pic-grid.component.html',
  styleUrls: ['./pic-grid.component.sass']
})
export class PicGridComponent implements OnInit {
  @Input() n: number = 0;

  elems: number[] = [];
  constructor() {
  }

  ngOnInit(): void {
    this.elems = []
    for (let i = 0; i < this.n; i++) {
      this.elems.push(i)
    }
    const o = new IntersectionObserver(console.log, {
      root: document
    })
  }
}
