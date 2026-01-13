export namespace main {
	
	export class DashboardSummary {
	    totalActiveProjects: number;
	    totalPlannedThisWeek: number;
	    totalActualThisWeek: number;
	    totalPlannedNextWeek: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalActiveProjects = source["totalActiveProjects"];
	        this.totalPlannedThisWeek = source["totalPlannedThisWeek"];
	        this.totalActualThisWeek = source["totalActualThisWeek"];
	        this.totalPlannedNextWeek = source["totalPlannedNextWeek"];
	    }
	}
	export class DashboardWeekData {
	    weekStartDate: string;
	    totalPlannedHours: number;
	    totalActualHours: number;
	    projectCount: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardWeekData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.weekStartDate = source["weekStartDate"];
	        this.totalPlannedHours = source["totalPlannedHours"];
	        this.totalActualHours = source["totalActualHours"];
	        this.projectCount = source["projectCount"];
	    }
	}

}

export namespace models {
	
	export class CreateGoalInput {
	    projectId: number;
	    title: string;
	    description: string;
	    targetDate: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateGoalInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.targetDate = source["targetDate"];
	    }
	}
	export class CreateMeetingInput {
	    projectId: number;
	    title: string;
	    meetingDate: string;
	    durationMinutes: number;
	    attendees: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateMeetingInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.meetingDate = source["meetingDate"];
	        this.durationMinutes = source["durationMinutes"];
	        this.attendees = source["attendees"];
	        this.notes = source["notes"];
	    }
	}
	export class CreateNoteInput {
	    projectId: number;
	    title: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateNoteInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.content = source["content"];
	    }
	}
	export class CreateProjectInput {
	    name: string;
	    totalSoldHours: number;
	    specialistHours: number;
	    startDate: string;
	    endDate: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.totalSoldHours = source["totalSoldHours"];
	        this.specialistHours = source["specialistHours"];
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	    }
	}
	export class Goal {
	    id: number;
	    projectId: number;
	    title: string;
	    description: string;
	    status: string;
	    targetDate: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Goal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.targetDate = source["targetDate"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GoalWithProject {
	    id: number;
	    projectId: number;
	    title: string;
	    description: string;
	    status: string;
	    targetDate: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    projectName: string;
	
	    static createFrom(source: any = {}) {
	        return new GoalWithProject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.targetDate = source["targetDate"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.projectName = source["projectName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Meeting {
	    id: number;
	    projectId: number;
	    title: string;
	    meetingDate: string;
	    durationMinutes: number;
	    attendees: string;
	    notes: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Meeting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.meetingDate = source["meetingDate"];
	        this.durationMinutes = source["durationMinutes"];
	        this.attendees = source["attendees"];
	        this.notes = source["notes"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MeetingWithProject {
	    id: number;
	    projectId: number;
	    title: string;
	    meetingDate: string;
	    durationMinutes: number;
	    attendees: string;
	    notes: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    projectName: string;
	
	    static createFrom(source: any = {}) {
	        return new MeetingWithProject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.meetingDate = source["meetingDate"];
	        this.durationMinutes = source["durationMinutes"];
	        this.attendees = source["attendees"];
	        this.notes = source["notes"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.projectName = source["projectName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Note {
	    id: number;
	    projectId: number;
	    title: string;
	    content: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Note(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NoteWithProject {
	    id: number;
	    projectId: number;
	    title: string;
	    content: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    projectName: string;
	
	    static createFrom(source: any = {}) {
	        return new NoteWithProject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.projectName = source["projectName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProjectWithStats {
	    id: number;
	    name: string;
	    totalSoldHours: number;
	    specialistHours: number;
	    startDate: string;
	    endDate: string;
	    isActive: boolean;
	    isPersistent: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    myHours: number;
	    totalWeeks: number;
	    totalPlannedHours: number;
	    totalActualHours: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectWithStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.totalSoldHours = source["totalSoldHours"];
	        this.specialistHours = source["specialistHours"];
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	        this.isActive = source["isActive"];
	        this.isPersistent = source["isPersistent"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.myHours = source["myHours"];
	        this.totalWeeks = source["totalWeeks"];
	        this.totalPlannedHours = source["totalPlannedHours"];
	        this.totalActualHours = source["totalActualHours"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateActualHoursInput {
	    entryId: number;
	    actualHours: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateActualHoursInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entryId = source["entryId"];
	        this.actualHours = source["actualHours"];
	    }
	}
	export class UpdateGoalInput {
	    id: number;
	    title: string;
	    description: string;
	    status: string;
	    targetDate: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateGoalInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.targetDate = source["targetDate"];
	    }
	}
	export class UpdateMeetingInput {
	    id: number;
	    title: string;
	    meetingDate: string;
	    durationMinutes: number;
	    attendees: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateMeetingInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.meetingDate = source["meetingDate"];
	        this.durationMinutes = source["durationMinutes"];
	        this.attendees = source["attendees"];
	        this.notes = source["notes"];
	    }
	}
	export class UpdateNoteInput {
	    id: number;
	    title: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateNoteInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	    }
	}
	export class UpdateProjectInput {
	    id: number;
	    name: string;
	    totalSoldHours: number;
	    specialistHours: number;
	    startDate: string;
	    endDate: string;
	    isActive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProjectInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.totalSoldHours = source["totalSoldHours"];
	        this.specialistHours = source["specialistHours"];
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	        this.isActive = source["isActive"];
	    }
	}
	export class WeeklyEntry {
	    id: number;
	    projectId: number;
	    weekStartDate: string;
	    weekNumber: number;
	    plannedHours: number;
	    actualHours: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new WeeklyEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.weekStartDate = source["weekStartDate"];
	        this.weekNumber = source["weekNumber"];
	        this.plannedHours = source["plannedHours"];
	        this.actualHours = source["actualHours"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WeeklyEntryWithProject {
	    id: number;
	    projectId: number;
	    weekStartDate: string;
	    weekNumber: number;
	    plannedHours: number;
	    actualHours: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    variance: number;
	    status: string;
	    isPastWeek: boolean;
	    projectName: string;
	
	    static createFrom(source: any = {}) {
	        return new WeeklyEntryWithProject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.weekStartDate = source["weekStartDate"];
	        this.weekNumber = source["weekNumber"];
	        this.plannedHours = source["plannedHours"];
	        this.actualHours = source["actualHours"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.variance = source["variance"];
	        this.status = source["status"];
	        this.isPastWeek = source["isPastWeek"];
	        this.projectName = source["projectName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WeeklyEntryWithStatus {
	    id: number;
	    projectId: number;
	    weekStartDate: string;
	    weekNumber: number;
	    plannedHours: number;
	    actualHours: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    variance: number;
	    status: string;
	    isPastWeek: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WeeklyEntryWithStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.weekStartDate = source["weekStartDate"];
	        this.weekNumber = source["weekNumber"];
	        this.plannedHours = source["plannedHours"];
	        this.actualHours = source["actualHours"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.variance = source["variance"];
	        this.status = source["status"];
	        this.isPastWeek = source["isPastWeek"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

