import { Component, ElementRef, HostListener, Input, OnInit, ViewChild } from '@angular/core';

type ImgData = {
  id: string,
  name: string,
  type: string,
  hash: string,
}

@Component({
  selector: 'app-pic-grid',
  templateUrl: './pic-grid.component.html',
  styleUrls: ['./pic-grid.component.sass']
})
export class PicGridComponent implements OnInit {
  imgs: ImgData[] = [];
  sel: number | undefined;

  constructor() {
  }

  ngOnInit(): void {
    this.fetchImages();
  }

  fetchImages() {
    fetch('/api/files').then(res => res.json()).then(res => {
      this.imgs = res;
    })
  }

  onImgClick(id: string) {
    this.sel = this.imgs.findIndex(img => img.id === id);
  }

  @HostListener("drop", ["$event"]) onDrop(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    if (!event.dataTransfer) return;
    const files = event.dataTransfer.files;
    const formData = new FormData();
    for (let i = 0; i < files.length; i++) {
      formData.append('files', files[i]);
    }
    console.log(event)
    fetch('/api/files', {
      method: 'POST',
      body: formData
    }).then(res => res.json()).then((res: ImgData[]) => {
      this.imgs = [...this.imgs, ...res.filter(e => !this.imgs.find(e2 => e2.hash === e.hash))];
    })
  }

  @HostListener("dragover", ["$event"]) onDragOver(event: DragEvent) {
    console.log('dragover');
    event.preventDefault();
  }
}
