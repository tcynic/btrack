import { useCallback, useState } from "react";
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
import { useQuery } from "./useQuery";

export function useDashboard() {
  const { dispatch } = useAppContext();
  const [dashboardParams, setDashboardParams] = useState({ weeksBack: 2, weeksForward: 4 });

  const summaryQuery = useQuery<DashboardSummary>({
    queryFn: GetDashboardSummary,
    onSuccess: (data) => {
      dispatch({ type: "SET_DASHBOARD_SUMMARY", payload: data });
    },
  });

  const weekDataQuery = useQuery<DashboardWeekData[]>({
    queryFn: useCallback(
      () => GetDashboardData(dashboardParams.weeksBack, dashboardParams.weeksForward),
      [dashboardParams]
    ),
    initialData: [],
    onSuccess: (data) => {
      dispatch({ type: "SET_DASHBOARD_WEEK_DATA", payload: data });
    },
  });

  const todayMeetingsQuery = useQuery<MeetingWithProject[]>({
    queryFn: useCallback(async () => {
      const today = new Date().toISOString().split("T")[0];
      const meetings = await GetMeetingsByDate(today);
      return (meetings || []) as MeetingWithProject[];
    }, []),
    initialData: [],
    onSuccess: (data) => {
      dispatch({ type: "SET_TODAY_MEETINGS", payload: data });
    },
  });

  const loadDashboardSummary = useCallback(async () => {
    await summaryQuery.refetch();
  }, [summaryQuery.refetch]);

  const loadDashboardData = useCallback(
    async (weeksBack: number = 2, weeksForward: number = 4) => {
      setDashboardParams({ weeksBack, weeksForward });
    },
    []
  );

  const loadTodayMeetings = useCallback(async () => {
    await todayMeetingsQuery.refetch();
  }, [todayMeetingsQuery.refetch]);

  const refreshDashboard = useCallback(async () => {
    await Promise.all([
      summaryQuery.refetch(),
      weekDataQuery.refetch(),
      todayMeetingsQuery.refetch(),
    ]);
  }, [summaryQuery.refetch, weekDataQuery.refetch, todayMeetingsQuery.refetch]);

  const isLoading = summaryQuery.isLoading || weekDataQuery.isLoading || todayMeetingsQuery.isLoading;
  const error = summaryQuery.error || weekDataQuery.error || todayMeetingsQuery.error;

  return {
    summary: summaryQuery.data || null,
    weekData: weekDataQuery.data || [],
    todayMeetings: todayMeetingsQuery.data || [],
    isLoading,
    error,
    loadDashboardSummary,
    loadDashboardData,
    loadTodayMeetings,
    refreshDashboard,
  };
}
