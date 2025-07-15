export interface User {
    id: string;
    name: string;
    email: string;
    lat: number;
    long: number;
    h3_id: string;
}

export interface NearbyResponse {
    users: User[];
    total: number;
}