import { createContext, useContext, useReducer, ReactNode } from "react";
import type {
  ProjectWithStats,
  DashboardSummary,
  DashboardWeekData,
  MeetingWithProject,
} from "../types";

interface AppState {
  projects: ProjectWithStats[];
  selectedProjectId: number | null;
  dashboardSummary: DashboardSummary | null;
  dashboardWeekData: DashboardWeekData[];
  todayMeetings: MeetingWithProject[];
  isLoading: boolean;
  error: string | null;
}

type AppAction =
  | { type: "SET_PROJECTS"; payload: ProjectWithStats[] }
  | { type: "ADD_PROJECT"; payload: ProjectWithStats }
  | { type: "UPDATE_PROJECT"; payload: ProjectWithStats }
  | { type: "DELETE_PROJECT"; payload: number }
  | { type: "SELECT_PROJECT"; payload: number | null }
  | { type: "SET_DASHBOARD_SUMMARY"; payload: DashboardSummary }
  | { type: "SET_DASHBOARD_WEEK_DATA"; payload: DashboardWeekData[] }
  | { type: "SET_TODAY_MEETINGS"; payload: MeetingWithProject[] }
  | { type: "SET_LOADING"; payload: boolean }
  | { type: "SET_ERROR"; payload: string | null };

const initialState: AppState = {
  projects: [],
  selectedProjectId: null,
  dashboardSummary: null,
  dashboardWeekData: [],
  todayMeetings: [],
  isLoading: false,
  error: null,
};

function appReducer(state: AppState, action: AppAction): AppState {
  switch (action.type) {
    case "SET_PROJECTS":
      return { ...state, projects: action.payload };
    case "ADD_PROJECT":
      return { ...state, projects: [action.payload, ...state.projects] };
    case "UPDATE_PROJECT":
      return {
        ...state,
        projects: state.projects.map((p) =>
          p.id === action.payload.id ? action.payload : p,
        ),
      };
    case "DELETE_PROJECT":
      return {
        ...state,
        projects: state.projects.filter((p) => p.id !== action.payload),
        selectedProjectId:
          state.selectedProjectId === action.payload
            ? null
            : state.selectedProjectId,
      };
    case "SELECT_PROJECT":
      return { ...state, selectedProjectId: action.payload };
    case "SET_DASHBOARD_SUMMARY":
      return { ...state, dashboardSummary: action.payload };
    case "SET_DASHBOARD_WEEK_DATA":
      return { ...state, dashboardWeekData: action.payload };
    case "SET_TODAY_MEETINGS":
      return { ...state, todayMeetings: action.payload };
    case "SET_LOADING":
      return { ...state, isLoading: action.payload };
    case "SET_ERROR":
      return { ...state, error: action.payload };
    default:
      return state;
  }
}

interface AppContextValue {
  state: AppState;
  dispatch: React.Dispatch<AppAction>;
}

const AppContext = createContext<AppContextValue | undefined>(undefined);

export function AppProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(appReducer, initialState);

  return (
    <AppContext.Provider value={{ state, dispatch }}>
      {children}
    </AppContext.Provider>
  );
}

export function useAppContext() {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error("useAppContext must be used within an AppProvider");
  }
  return context;
}
