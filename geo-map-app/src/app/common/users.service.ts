import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { NearbyResponse } from './interfaces/user.interface';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  constructor(private readonly http: HttpClient) { }

  loadNearbyUsers(lat: number, lng: number): Observable<NearbyResponse> {
    const url = `${environment.apiBaseUrl}/api/nearby`
    return this.http.post<NearbyResponse>(url, {lat, long: lng});    
  }
}
