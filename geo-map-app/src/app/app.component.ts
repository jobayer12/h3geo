import { Component, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import * as L from 'leaflet';

interface User {
  id: string;
  name: string;
  email: string;
  lat: number;
  long: number;
  h3_id: string;
}

interface NearbyResponse {
  users: User[];
  total: number;
}

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, HttpClientModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements AfterViewInit {
  private map!: L.Map;
  private clickMarker?: L.CircleMarker; 
  private clickCircle?: L.Circle;
  private userMarkers: L.Marker[] = [];
  users: User[] = [];
  loading = false;

  constructor(private http: HttpClient) {}

  ngAfterViewInit(): void {
    this.initMap();
  }

  private initMap(): void {
    this.map = L.map('map').setView([23.685, 90.3563], 7);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 18,
      attribution: '© OpenStreetMap contributors'
    }).addTo(this.map);

    this.map.on('click', (e: L.LeafletMouseEvent) => this.onMapClick(e));
  }

  private onMapClick(e: L.LeafletMouseEvent): void {
    const { lat, lng } = e.latlng;

    // Remove previous click marker and circle
    this.clearPreviousClickElements();

    // Add new click marker
    // this.clickMarker = L.circleMarker([lat, lng], {
    //   radius: 8,
    //   fillColor: 'red',
    //   color: 'red',
    //   weight: 2,
    //   opacity: 1,
    //   fillOpacity: 0.8
    // }).addTo(this.map);
    // Add circle with 5.2km radius
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
    if (this.clickMarker) {
      this.map.removeLayer(this.clickMarker);
      this.clickMarker = undefined;
    }
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
    this.loading = true;
    this.users = [];
    this.clearUserMarkers();

    this.http.post<NearbyResponse>('http://localhost:8080/api/nearby', { lat, long: lng })
      .subscribe({
        next: (response) => {
          this.users = response.users;
          this.addUserMarkers(response.users);
          this.loading = false;
        },
        error: (error) => {
          console.error('Error fetching nearby users:', error);
          this.loading = false;
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