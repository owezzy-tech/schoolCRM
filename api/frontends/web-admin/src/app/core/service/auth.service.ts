import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, map, of } from 'rxjs';
import { User } from '../models/user';
import { Role } from '@core/models/role';
import { environment } from 'environments/environment';

interface LoginResponse {
  accessToken: string;
  tokenType: string;
  expiresAt: string;
  expiresIn: number;
  user: {
    id: string;
    name: string;
    email: string;
    roles: Role[];
  };
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly emptyUser = {} as User;
  private currentUserSubject: BehaviorSubject<User>;
  public currentUser: Observable<User>;

  constructor(private http: HttpClient) {
    this.currentUserSubject = new BehaviorSubject<User>(
      JSON.parse(localStorage.getItem('currentUser') || 'null') ?? this.emptyUser
    );
    this.currentUser = this.currentUserSubject.asObservable();
  }

  public get currentUserValue(): User {
    return this.currentUserSubject.value;
  }

  login(username: string, password: string): Observable<User> {
    return this.http
      .post<LoginResponse>(`${environment.apiUrl}/auth/login`, {
        email: username,
        password,
      })
      .pipe(
        map((response) => {
          const user = this.toUser(response);
          localStorage.setItem('currentUser', JSON.stringify(user));
          this.currentUserSubject.next(user);
          return user;
        })
      );
  }

  logout() {
    // remove user from local storage to log user out
    localStorage.removeItem('currentUser');
    this.currentUserSubject.next(this.emptyUser);
    return of({ success: false });
  }

  private toUser(response: LoginResponse): User {
    const [firstName, ...rest] = response.user.name.split(' ');
    const roles = response.user.roles ?? [];

    return {
      id: response.user.id,
      img: this.imageForRole(roles[0]),
      username: response.user.email,
      firstName: firstName || response.user.name,
      lastName: rest.join(' '),
      role: roles[0] ?? Role.Student,
      roles,
      token: response.accessToken,
      expiresAt: response.expiresAt,
    };
  }

  private imageForRole(role: Role | undefined): string {
    switch (role) {
      case Role.SuperAdmin:
      case Role.SchoolAdmin:
        return 'assets/images/user/admin.jpg';
      case Role.Teacher:
        return 'assets/images/user/teacher.jpg';
      case Role.Student:
        return 'assets/images/user/student.jpg';
      default:
        return 'assets/images/user/user.jpg';
    }
  }
}
