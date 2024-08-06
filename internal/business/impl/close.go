package impl

func (g *gophermart) Close() {
	g.wgIn.Wait()
	close(g.in)

	g.accrualSvc.Close()

	g.wgOut.Wait()

	g.logger.Info("goodbye from business-svc")
}
