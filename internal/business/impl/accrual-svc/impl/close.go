package impl

func (s *svc) Close() {
	s.wgWorkers.Wait()
	s.wgProcess.Wait()

	close(s.outUpdated)
	close(s.outNotUpdated)

	s.logger.Info("goodbye from accrual-svc")
}
