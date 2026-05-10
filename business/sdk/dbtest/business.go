package dbtest

import (
	"time"

	"github.com/owezzy/schoolCRM/business/domain/auditbus"
	"github.com/owezzy/schoolCRM/business/domain/auditbus/extensions/auditotel"
	"github.com/owezzy/schoolCRM/business/domain/auditbus/stores/auditdb"
	"github.com/owezzy/schoolCRM/business/domain/homebus"
	"github.com/owezzy/schoolCRM/business/domain/homebus/extensions/homeotel"
	"github.com/owezzy/schoolCRM/business/domain/homebus/stores/homedb"
	"github.com/owezzy/schoolCRM/business/domain/productbus"
	"github.com/owezzy/schoolCRM/business/domain/productbus/extensions/productotel"
	"github.com/owezzy/schoolCRM/business/domain/productbus/stores/productdb"
	"github.com/owezzy/schoolCRM/business/domain/userbus"
	"github.com/owezzy/schoolCRM/business/domain/userbus/extensions/useraudit"
	"github.com/owezzy/schoolCRM/business/domain/userbus/extensions/userotel"
	"github.com/owezzy/schoolCRM/business/domain/userbus/stores/usercache"
	"github.com/owezzy/schoolCRM/business/domain/userbus/stores/userdb"
	"github.com/owezzy/schoolCRM/business/domain/vproductbus"
	"github.com/owezzy/schoolCRM/business/domain/vproductbus/extensions/vproductotel"
	"github.com/owezzy/schoolCRM/business/domain/vproductbus/stores/vproductdb"
	"github.com/owezzy/schoolCRM/business/sdk/delegate"
	"github.com/owezzy/schoolCRM/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	Audit    auditbus.ExtBusiness
	Home     homebus.ExtBusiness
	Product  productbus.ExtBusiness
	User     userbus.ExtBusiness
	VProduct vproductbus.ExtBusiness
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)

	auditOtelExt := auditotel.NewExtension()
	auditStorage := auditdb.NewStore(log, db)
	auditBus := auditbus.NewBusiness(log, auditStorage, auditOtelExt)

	userOtelExt := userotel.NewExtension()
	userAuditExt := useraudit.NewExtension(auditBus)
	userStorage := usercache.NewStore(log, userdb.NewStore(log, db), time.Hour)
	userBus := userbus.NewBusiness(log, delegate, userStorage, userOtelExt, userAuditExt)

	productOtelExt := productotel.NewExtension()
	productStorage := productdb.NewStore(log, db)
	productBus := productbus.NewBusiness(log, userBus, delegate, productStorage, productOtelExt)

	homeOtelExt := homeotel.NewExtension()
	homeStorage := homedb.NewStore(log, db)
	homeBus := homebus.NewBusiness(log, userBus, delegate, homeStorage, homeOtelExt)

	vproductOtelExt := vproductotel.NewExtension()
	vproductStorage := vproductdb.NewStore(log, db)
	vproductBus := vproductbus.NewBusiness(vproductStorage, vproductOtelExt)

	return BusDomain{
		Delegate: delegate,
		Audit:    auditBus,
		Home:     homeBus,
		Product:  productBus,
		User:     userBus,
		VProduct: vproductBus,
	}
}
