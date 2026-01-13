import { useCallback } from "react";
import { useAppContext } from "../context/AppContext";
import type {
  DashboardSummary,
  DashboardWeekData,
  MeetingWithProject,
} from "../types";
import {
  GetDashboardSummary,
  GetDashboardData,
  GetMeetingsByDate,
} from "../../wailsjs/go/main/App";

export function useDashboard() {
  const { state, dispatch } = useAppContext();

  const loadDashboardSummary = useCallback(async () => {
    dispatch({ type: "SET_LOADING", payload: true });
    dispatch({ type: "SET_ERROR", payload: null });
    try {
      const summary = await GetDashboardSummary();
      dispatch({
        type: "SET_DASHBOARD_SUMMARY",
        payload: summary as DashboardSummary,
      });
    } catch (err) {
      dispatch({ type: "SET_ERROR", payload: String(err) });
    } finally {
      dispatch({ type: "SET_LOADING", payload: false });
    }
  }, [dispatch]);

  const loadDashboardData = useCallback(
    async (weeksBack: number = 4, weeksForward: number = 12) => {
      dispatch({ type: "SET_LOADING", payload: true });
      dispatch({ type: "SET_ERROR", payload: null });
      try {
        const data = await GetDashboardData(weeksBack, weeksForward);
        dispatch({
          type: "SET_DASHBOARD_WEEK_DATA",
          payload: data as DashboardWeekData[],
        });
      } catch (err) {
        dispatch({ type: "SET_ERROR", payload: String(err) });
      } finally {
        dispatch({ type: "SET_LOADING", payload: false });
      }
    },
    [dispatch],
  );

  const loadTodayMeetings = useCallback(async () => {
    try {
      const today = new Date().toISOString().split("T")[0];
      const meetings = await GetMeetingsByDate(today);
      dispatch({
        type: "SET_TODAY_MEETINGS",
        payload: (meetings || []) as MeetingWithProject[],
      });
    } catch (err) {
      console.error("Failed to load today meetings:", err);
    }
  }, [dispatch]);

  const refreshDashboard = useCallback(async () => {
    await Promise.all([
      loadDashboardSummary(),
      loadDashboardData(),
      loadTodayMeetings(),
    ]);
  }, [loadDashboardSummary, loadDashboardData, loadTodayMeetings]);

  return {
    summary: state.dashboardSummary,
    weekData: state.dashboardWeekData,
    todayMeetings: state.todayMeetings,
    isLoading: state.isLoading,
    error: state.error,
    loadDashboardSummary,
    loadDashboardData,
    loadTodayMeetings,
    refreshDashboard,
  };
}
