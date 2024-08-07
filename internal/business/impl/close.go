package impl

func (g *gophermart) Close() {
	g.killIn = nil
	g.wgIn.Wait()
	close(g.in)
	close(g.errChan)

	g.accrualSvc.Close()

	g.wgOut.Wait()

	g.logger.Info("goodbye from business-svc")
}
