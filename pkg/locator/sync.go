package locator

import "context"

func (p *PDBLocator) waitForCacheSync(ctx context.Context) error {
	p.infFactory.WaitForCacheSync(ctx.Done())
	return ctx.Err()
}
