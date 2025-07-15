import { Component, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as L from 'leaflet';
import { UsersService } from './common/users.service';
import { User } from './common/interfaces/user.interface';
import { Subject, takeUntil } from 'rxjs';


@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule],
  providers: [UsersService],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements AfterViewInit {
  private map!: L.Map;
  private clickCircle?: L.Circle;
  private userMarkers: L.Marker[] = [];
  users: User[] = [];
  cancelPreviousRequest = new Subject();

  constructor(
    private readonly userService: UsersService
  ) {}

  ngAfterViewInit(): void {
    this.initMap();
  }

  private initMap(): void {
    this.map = L.map('map').setView([23.685, 90.3563], 2);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 18,
      attribution: '© OpenStreetMap contributors'
    }).addTo(this.map);

    this.map.on('click', (e: L.LeafletMouseEvent) => this.onMapClick(e));
  }

  private onMapClick(e: L.LeafletMouseEvent): void {
    const { lat, lng } = e.latlng;
    this.clearPreviousClickElements();

    this.clickCircle = L.circle([lat, lng], {
      color: 'red',
      fillColor: '#f03',
      fillOpacity: 0.2,
      radius: 5500
    }).addTo(this.map);

    this.map.setView([lat, lng], 12, {
      animate: true,
      duration: 0.5
    });


    // Fetch nearby users
    this.getNearbyUsers(lat, lng);
  }

  private clearPreviousClickElements(): void {
    if (this.clickCircle) {
      this.map.removeLayer(this.clickCircle);
      this.clickCircle = undefined;
    }
  }

  private clearUserMarkers(): void {
    this.userMarkers.forEach(marker => {
      this.map.removeLayer(marker);
    });
    this.userMarkers = [];
  }

  private getNearbyUsers(lat: number, lng: number): void {
    this.cancelPreviousRequest.next(null);
    this.cancelPreviousRequest.complete(); 
    this.cancelPreviousRequest = new Subject(); 
    this.users = [];
    this.clearUserMarkers();

    this.userService.loadNearbyUsers(lat, lng).pipe(takeUntil(this.cancelPreviousRequest))
    .subscribe({
      next: (response) => {
        this.users = response.users;
        this.addUserMarkers(response.users);
      },
      error: (error) => {
        console.error('Error fetching nearby users:', error);
      }
    });
  }

  private addUserMarkers(users: User[]): void {
    (users || []).forEach(user => {
      const userMarker = L.marker([user.lat, user.long], {
        icon: L.icon({
          iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
          shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
          iconSize: [20, 32],
          iconAnchor: [10, 32],
          popupAnchor: [1, -34],
          shadowSize: [32, 32]
        })
      }).addTo(this.map);

      // Add popup with user info
      userMarker.bindPopup(`
        <div style="text-align: center;">
          <strong>${user.name}</strong>
        </div>
      `);

      this.userMarkers.push(userMarker);
    });
  }
}